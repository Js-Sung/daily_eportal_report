# 每日保平安
This project handles the daily report of healthy of UESTC.  
![Image](https://github.com/Js-Sung/daily_eportal_report/blob/master/useless/signed.png)  

## Dependency
golang, chromedp, chrome browser(Or chrome-based browsers. I am using Cent Browser that works fine:)).

## Feature
1. This Project handles daily reports of healthy automatically.
2. A screenshot will be saved if scripts work successfully.
3. Scripts will quit if running time exceeds 3 minutes.
4. Scripts won't retry if anything wrong is encountered.  


## HOWTO
1. Install golang, chromedp, chrome browser.
2. Download all files in this repository. Modify the line 6-8 of main.cmd.
3. Launch main.cmd to see if it works well.
4. Add main.cmd as a daily scheduled task. Add two triggers just in case it fails the first time.  
![Image](https://github.com/Js-Sung/daily_eportal_report/blob/master/useless/1.png)  

## Reference
1. [bypass headless chrome detection
(code by kenshaw)](https://github.com/chromedp/chromedp/issues/396#issuecomment-503351342)
2. [IT IS *NOT* POSSIBLE TO DETECT AND BLOCK CHROME HEADLESS](https://intoli.com/blog/not-possible-to-block-chrome-headless/)
