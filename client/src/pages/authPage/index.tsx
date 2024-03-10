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
            
            if (resp.data.redirectUrl) {
                redirectUrl = resp.data.redirectUrl.includes(window.location.origin) ? resp.data.redirectUrl : `${window.location.origin}/#${routes.redirectPage.path}`;
            } 
                // 如果 resp.data.redirectUrl 为空，则从当前页面链接获取 redirect_url 参数
                const urlParams = new URLSearchParams(window.location.search);
                console.log("window.location---", window.location, urlParams)
                const urlRedirectUrl = urlParams.get('redirect_url');
                if (urlRedirectUrl) {
                    redirectUrl = urlRedirectUrl;
                    console.log("浏览器获取的重定向链接", redirectUrl)
                }
            if (resp1.code === 0) {
                setTimeout(() => {
                    window.location.href = `https://open.weixin.qq.com/wxaopen/safe/bindcomponent?component_appid=${resp.data.appid}&pre_auth_code=${resp1.data.preAuthCode}&auth_type=6&redirect_uri=${encodeURIComponent(redirectUrl)}#wechat_redirect`
                }, 10000);
            }
        }
    }

    return (
        <div />
    )
}
