import {type UseFetchOptions, useFetch} from "nuxt/app";

type Methods = "GET" | "POST" | "DELETE" | "PUT";

export interface IPageData<T> {
    pageNo: number;
    pageSize: number;
    totalPage: number;
    totalCount: number;
    list: T[];
}

export interface IListData<T> {
    list: T[];
}

export interface IBaseResponse {
    code: number;
    message: string;
}

export interface IResponse<T> {
    code: number;
    data?: T;
    message: string;
}

class HttpRequest {
    request<T = any>(
        url: string,
        method: Methods,
        data: any,
        options?: UseFetchOptions<T>,
    ) {
        return new Promise((resolve, reject) => {
            const newOptions: UseFetchOptions<T> = {
                baseURL: url,
                method: method,
                ...options,
            };

            if (method === "GET" || method === "DELETE") {
                newOptions.params = data;
            }
            if (method === "POST" || method === "PUT") {
                newOptions.body = data;
            }
            useFetch(url, newOptions)
                .then((res) => {
                    resolve(res);
                })
                .catch((error) => {
                    reject(error);
                });
        });
    }


    // 封装常用方法

    get<T = any>(url: string, params?: any, options?: UseFetchOptions<T>) {
        return this.request(url, "GET", params, options);
    }

    post<T = any>(url: string, data: any, options?: UseFetchOptions<T>) {
        return this.request(url, "POST", data, options);
    }

    put<T = any>(url: string, data: any, options?: UseFetchOptions<T>) {
        return this.request(url, "PUT", data, options);
    }

    delete<T = any>(url: string, params: any, options?: UseFetchOptions<T>) {
        return this.request(url, "DELETE", params, options);
    }
}

const httpRequest = new HttpRequest();

export default httpRequest