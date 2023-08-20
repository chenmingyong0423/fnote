/**
 * 函数防抖
 * @param {*} fn
 * @param {*} delay
 * @returns
 */
export function debounce(fn, delay) {
    delay = delay || 1000
    let timer = null
    return function () {
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const context = this
        // eslint-disable-next-line prefer-rest-params
        const arg = arguments
        if (timer)
            clearTimeout(timer)

        timer = setTimeout(() => {
            fn.apply(context, arg)
        }, delay)
    }
}
/**
 * 节流函数
 * @param {*} fn
 * @param {*} delay
 * @returns
 */
export function throttle(fn, delay = 300) {
    let timer = null
    return function () {
        // eslint-disable-next-line @typescript-eslint/no-this-alias
        const context = this
        // eslint-disable-next-line prefer-rest-params
        const args = arguments
        if (!timer) {
            timer = setTimeout(() => {
                fn.apply(context, args)
                clearTimeout(timer)
            }, delay)
        }
    }
}