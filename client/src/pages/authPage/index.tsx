import {useEffect} from "react";
import {request} from "../../utils/axios";
import {getComponentInfoRequest, getPreAuthCodeRequest} from "../../utils/apis";
import {routes} from "../../config/route";

export default function AuthPage() {

    useEffect(() => {
        jumpAuthPage()
    }, [])

    const jumpAuthPage = async () => {
        let redirectUrl = ''
        const resp = await request({
            request: getComponentInfoRequest,
            noNeedCheckLogin: true
        })
        if (resp.code === 0) {
            const resp1 = await request({
                request: getPreAuthCodeRequest,
                noNeedCheckLogin: true
            })
            
            if (resp.data.redirectUrl && resp.data.redirectUrl != "") {
                redirectUrl = resp.data.redirectUrl.includes(window.location.origin) ? resp.data.redirectUrl : `${window.location.origin}/#${routes.redirectPage.path}`;
                console.log("数据库的", redirectUrl, resp.data)
            } else {
                // 如果 resp.data.redirectUrl 为空，则从当前页面链接获取 redirect_url 参数
                const url = new URL(window.location.href);
                const redirectUri = url.hash.slice(25)
                if (redirectUri && redirectUri != '') {
                  redirectUrl = `${window.location.origin}/#${routes.redirectPage.path}?redirect_url=${redirectUri}`;
                  console.log("浏览器带的", redirectUrl)
                }            
            }
            if (resp1.code === 0) {
                setTimeout(() => {
                    window.location.href = `https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=${resp.data.appid}&pre_auth_code=${resp1.data.preAuthCode}&auth_type=6&redirect_uri=${encodeURIComponent(redirectUrl)}`
                }, 20000);
            }
        }
    }

    return (
        <div />
    )
}
