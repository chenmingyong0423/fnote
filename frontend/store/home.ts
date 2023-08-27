export const useHomeStore = defineStore("home", {
    state: () => ({
        myHeaderBg: 'bg-#000/20 backdrop-blur-20',
        headerTextColor: '#fff',
        menuList: {
            "code": 200,
            "message": "OK",
            "data": {
                "list": [
                    {
                        "name": "后端",
                        "route": "/category/backend"
                    },
                    {
                        "name": "前端",
                        "route": "/category/frontend"
                    }
                ]
            }
        }
    })
});