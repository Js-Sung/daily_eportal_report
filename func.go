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
)



var exe_path, temp_path, id, passwd string



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

	// see: https://intoli.com/blog/not-possible-to-block-chrome-headless/
	const script = `(function(w, n, wn) {
		// Pass the Webdriver Test.
		Object.defineProperty(n, 'webdriver', {
			get: () => false,
		});

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
		return wn.permissions.query = (parameters) => (
			parameters.name === 'notifications' ?
			Promise.resolve({ state: Notification.permission }) :
			originalQuery(parameters)
		);

	})(window, navigator, window.navigator);`
	

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),
		chromedp.Flag("window-size", "1024,800"),
		chromedp.Flag("enable-automation", false),
		//chromedp.ProxyServer("http://192.168.1.217:8088"),
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

	// run
	var b1 []byte
	var scriptID page.ScriptIdentifier

	if err := chromedp.Run(taskCtx,
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
		chromedp.WaitVisible(`#captcha-id`, chromedp.ByID),
		
		// try to sign in...
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
		
		chromedp.WaitVisible(`#widget-recommendAndNew-01 > div.clearfix.active.card-recommend-new-main.style-scope.pc-card-html-4786696181714491-01 > widget-app-item:nth-child(1) > div > div > div.widget-information.style-scope.pc-card-html-4786696181714491-01 > div`, chromedp.ByID),
		
		chromedp.Navigate(`http://eportal.uestc.edu.cn/qljfwapp/sys/lwReportEpidemicStu/index.do`),
		// click on the "add" button
		chromedp.Click(`body > main > article > section > div.bh-mb-16 > div`, chromedp.NodeReady),
		// wait for the table visible
		chromedp.WaitVisible(`#emapForm > div > div:nth-child(3)`, chromedp.ByID),
		// click on the "save" button
		chromedp.Click(`#save`, chromedp.NodeReady),
		// click on the "sure" button
		chromedp.Click(`a.bh-dialog-btn.bh-bg-primary.bh-color-primary-5`, chromedp.NodeReady),
		// take screenshot
		chromedp.CaptureScreenshot(&b1),
	); err != nil {
		log.Fatal(err)
	}
	
	// save screenshot to a file
	if err := ioutil.WriteFile(id+"-log.png", b1, 0644); err != nil {
		log.Fatal(err)
	}
	
	
}
