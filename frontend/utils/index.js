/**
 * 函数防抖
 * @param {*} fn
 * @param {*} delay
 * @returns
 */
export function debounce(fn, delay = 300) {
    if (timer != null) {
        clearTimeout(timer)
        timer = null
    }
    timer = setTimeout(fn, delay)
}
/**
 * 节流函数
 * @param {*} fn
 * @param {*} delay
 * @returns
 */


export function throttle(fn, delay = 300) {
    let throttleTimer = null
    return function (...args) {
        if (throttleTimer == null) {
            throttleTimer = setTimeout(() => {
                fn.call(this, ...args)
                clearTimeout(throttleTimer)
                throttleTimer = null
            }, delay);
        }
    }
}