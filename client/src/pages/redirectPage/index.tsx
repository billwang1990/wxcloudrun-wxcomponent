import {useEffect} from "react";
import {request} from "../../utils/axios";
import {getComponentInfoRequest} from "../../utils/apis";
import {routes} from "../../config/route";

export default function RedirectPage() {

    useEffect(() => {
        jumpRealPage()
    }, [])

    const jumpRealPage = async () => {
        const resp = await request({
            request: getComponentInfoRequest,
            noNeedCheckLogin: true
        })
        if (resp.code === 0) {
            const originalUrl = window.location.href
            if (originalUrl.includes("redirect_url")) {
                // 解析原始URL以获取hash部分
                const url = new URL(originalUrl);
                const hashParams = new URLSearchParams(url.hash.split('?')[1]); // 获取hash中的查询参数部分

                // 提取redirect_url的值，并从参数中移除它
                const redirectUrl = hashParams.get('redirect_url');
                hashParams.delete('redirect_url');

                // 将剩余的所有参数附加到redirectUrl上
                const finalUrl = `${redirectUrl}?${hashParams.toString()}`;

                console.log(finalUrl);
                window.location.href = finalUrl;
            } else {
                window.location.href = resp.data.redirectUrl + window.location.hash.replaceAll(`#${routes.redirectPage.path}`, '')
            }
        }
    }

    return (
        <div />
    )
}
