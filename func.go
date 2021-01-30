package main

import (
	"flag"
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/cdproto/cdp"
)



var exe_path, temp_path, id, passwd string

// run
var b1 []byte
var scriptID page.ScriptIdentifier
var nodes []*cdp.Node

// see: https://intoli.com/blog/not-possible-to-block-chrome-headless/
// see: https://github.com/paulirish/headless-cat-n-mouse/
const script = `(function(w, n, wn) {
	// Pass the Webdriver Test.
	const newProto = n.__proto__;
	delete newProto.webdriver;
	n.__proto__ = newProto;

	// Pass the Plugins Length Test.
	// Overwrite the plugins property to use a custom getter.
	Object.defineProperty(n, 'plugins', {
		// This just needs to have length > 0 for the current test,
		// but we could mock the plugins too if necessary.
		get: () => [1, 2, 3, 4, 5],
	});

	// Pass the Languages Test.
	// Overwrite the plugins property to use a custom getter.
	Object.defineProperty(n, 'languages', {
		get: () => ['zh-CN', 'zh'],
	});

	// Pass the Chrome Test.
	// We can mock this in as much depth as we need for the test.
	w.chrome = {
		runtime: {},
	};

	// Pass the Permissions Test.
	const originalQuery = wn.permissions.query;
	wn.permissions.__proto__.query = parameters =>
	parameters.name === 'notifications' ? Promise.resolve({state: Notification.permission}) : originalQuery(parameters);

	// Inspired by: https://github.com/ikarienator/phantomjs_hide_and_seek/blob/master/5.spoofFunctionBind.js
	const oldCall = Function.prototype.call;
	function call() {
		return oldCall.apply(this, arguments);
	}
	Function.prototype.call = call;

	const nativeToStringFunctionString = Error.toString().replace(/Error/g, "toString");
	const oldToString = Function.prototype.toString;

	function functionToString() {
		if (this === wn.permissions.query) {
			return "function query() { [native code] }";
		}
		if (this === functionToString) {
			return nativeToStringFunctionString;
		}
		return oldCall.call(oldToString, this);
	}
	Function.prototype.toString = functionToString;

})(window, navigator, window.navigator);`

func main() {

	flag.StringVar(&id, "id", "201821019876", "your id number")
	flag.StringVar(&passwd, "passwd", "123456", "your password")
	flag.StringVar(&exe_path, "exe_path", `D:\Program Files\CentBrowser\Application\chrome.exe`, "your chrome path")
	flag.StringVar(&temp_path, "temp_path", `D:\temp`, "your temp folder for chrome")
	flag.Parse()
	fmt.Println("id="+id, "password="+passwd, "exe path="+exe_path, "temp path="+temp_path)

	
	dir, err := ioutil.TempDir(temp_path, "chromedp-example")
	if err != nil {
		panic(err)
	}else{
		fmt.Println("temp dir = "+dir)
	}
	defer os.RemoveAll(dir)


	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),
		chromedp.Flag("window-size", "1024,800"),
		chromedp.Flag("enable-automation", false),
		//chromedp.ProxyServer("http://localhost:8888"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36`),

		chromedp.UserDataDir(dir),
		chromedp.ExecPath(exe_path),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// create a timeout
	taskCtx, cancel = context.WithTimeout(taskCtx, 180*time.Second)
	defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}
	
	// navigate to eportal, input usrname passwd, click on 'sign in' button
	fmt.Println("try to sign in.")
	if err := chromedp.Run(taskCtx, task1()); err != nil {
		log.Fatal(err)
	}
	
	// check if Captcha exists
	time.Sleep(500 * time.Millisecond)
	if err := chromedp.Run(taskCtx,
		chromedp.Nodes("#captcha-id", &nodes, chromedp.AtLeast(0)),
	); err != nil {
		log.Fatal(err)
	} else {
	
		if nodes != nil {
			fmt.Print("let's check if Captcha exists: ")
			if len(nodes) > 0 {
				fmt.Println("need to bypass Captcha.")
				if err := chromedp.Run(taskCtx, task3()); err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Println("no Captcha.")
			}
			fmt.Println("logged in.")
		} else {
			fmt.Println("warning: can not determine if Captcha exists.")
		}
		
		// click on 'add; button 
		if err := chromedp.Run(taskCtx,task2()); err != nil {
			log.Fatal(err)
		}
		
		// check if you have clocked in before
		time.Sleep(500 * time.Millisecond)
		if err := chromedp.Run(taskCtx,
			chromedp.Nodes("div.bh-pop.bh-card.bh-card-lv4.bh-dialog-con", &nodes, chromedp.AtLeast(0)),
		); err != nil {
			log.Fatal(err)
		} else {
		
			if nodes != nil {
				fmt.Print("let's check if you have clocked in before: ")
				if len(nodes) < 1 {
					fmt.Println("you need to clock in.")
					if err := chromedp.Run(taskCtx,
						// click on the "save" button
						chromedp.Click(`#save`, chromedp.NodeReady),
						// click on the "sure" button
						chromedp.Click(`a.bh-dialog-btn.bh-bg-primary.bh-color-primary-5`, chromedp.NodeReady),
					); err != nil {
						log.Fatal(err)
					}
				} else {
					fmt.Println("you have already clocked in.")
				}
				// take screenshot
				if err := chromedp.Run(taskCtx,
					chromedp.CaptureScreenshot(&b1),
				); err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Println("warning: can not determine if you have already clocked in.")
			}
		}
	
	}

	
	// save screenshot to a file
	if err := ioutil.WriteFile(id+"-log.png", b1, 0644); err != nil {
		log.Fatal(err)
	}
	
}


func task1() chromedp.Tasks{
	return chromedp.Tasks{
		network.Enable(),
		chromedp.Emulate(device.IPhone7),
		chromedp.EmulateViewport(1024, 800),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			scriptID, err = page.AddScriptToEvaluateOnNewDocument(script).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
		
		//chromedp.Navigate(`https://intoli.com/blog/not-possible-to-block-chrome-headless/chrome-headless-test.html`),
		chromedp.Navigate(`https://idas.uestc.edu.cn/authserver/login?service=http%3A%2F%2Feportal.uestc.edu.cn%2Flogin%3Fservice%3Dhttp%3A%2F%2Feportal.uestc.edu.cn%2Fnew%2Findex.html`),
		chromedp.WaitVisible(`#casLoginForm`, chromedp.ByID),
		chromedp.SetValue(`#mobileUsername`, id, chromedp.ByID),
		chromedp.SetValue(`#mobilePassword`, passwd, chromedp.ByID),
		chromedp.Click(`#load`, chromedp.NodeReady),	
		
	}
}

func task2() chromedp.Tasks{
	return chromedp.Tasks{
		chromedp.WaitVisible(`#widget-recommendAndNew-01`, chromedp.ByID),
		
		chromedp.Navigate(`http://eportal.uestc.edu.cn/qljfwapp/sys/lwReportEpidemicStu/index.do`),
		// click on the "add" button
		chromedp.Click(`body > main > article > section > div.bh-mb-16 > div`, chromedp.NodeReady),
	}
}


func task3() chromedp.Tasks{
	return chromedp.Tasks{
		// try to bypass Captcha...
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, exp, err := runtime.Evaluate(`document.querySelector("#casLoginForm").submit();`).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
	}
}
