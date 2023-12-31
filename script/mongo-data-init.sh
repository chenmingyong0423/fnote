mongosh -- "$MONGO_INITDB_DATABASE" <<EOF
db = db.getSiblingDB('$MONGO_INITDB_DATABASE')
db.auth('$MONGO_USERNAME', '$MONGO_PASSWORD');

// ----------------------------
// Collection structure for categories
// ----------------------------
db.getCollection("categories").drop();
db.createCollection("categories");
db.categories.createIndex({ "route": -1 });



// ----------------------------
// Collection structure for comment
// ----------------------------
db.getCollection("comment").drop();
db.createCollection("comment");
db.comment.createIndex({ "post_info.post_id": 1 });
db.comment.createIndex({ "status": -1 });

// ----------------------------
// Collection structure for configs
// ----------------------------
db.getCollection("configs").drop();
db.createCollection("configs");

db.configs.createIndex({ "typ": 1 });

// ----------------------------
// Documents of configs
// ----------------------------
// 网站信息
db.getCollection("configs").insertOne({
    _id: "webmaster",
    typ: "webmaster",
    props: {
        name: "fnote",
        postCount: 0,
        categoryCount: 0,
        websiteViews: 0,
        websiteLiveTime: Date.now(),
        profile: "请及时到后台修改网站的信息以及相关 seo 等配置，以便正常使用。",
        picture: "/fnote_logo.png",
        websiteIcon: "/fnote_logo.pne",
        domain: ""
    },
    create_time: Date.now(),
    update_time: Date.now()
});
// seo 配置
db.getCollection("configs").insertOne({
    _id: "webmaster",
    "typ": "seo meta",
    "props": {
        "title": "fnote",
        "ogTitle": "fnote",
        "description": "fnote",
        "ogImage": "",
        "twitterCard": "",
        "baidu-site-verification": "",
        "keywords": "fnote,blog,BLOG",
        "author": "fnote",
        "robots": "fnote,blog"
    },
    create_time: Date.now(),
    update_time: Date.now()
});
// 评论开关配置
db.getCollection("configs").insertOne({
    _id: "comment",
    typ: "comment",
    props: {
        status: true
    },
    create_time: Date.now(),
    update_time: Date.now()
});
// 友链开关配置
db.getCollection("configs").insertOne({
    _id: "friend",
    typ: "friend",
    props: {
        status: false
    },
    create_time: Date.now(),
    update_time: Date.now()
});
// 邮件配置
db.getCollection("configs").insertOne({
    _id: "emailConfig",
    "typ": "emailConfig",
    "props": {
        "host": "",
        "port": ,
        "account": "",
        "password": "",
        "email": ""
    },
    create_time: Date.now(),
    update_time: Date.now()
});
// 公告配置
db.getCollection("configs").insertOne({
    _id: "notice",
    typ: "notice",
    "props": {
        "title": "暂无最新公告",
        "content": "暂无最新公告",
        "publish_time": Date.now(),
        "enabled": true
    },
    create_time: Date.now(),
    update_time: Date.now()
});
// 社交外链配置
db.getCollection("configs").insertOne({
    _id: "social",
    typ: "social",
    "props": {
        "social_info_list": [
            {
                "social_name": "fnote",
                "social_value": "https://github.com/chenmingyong0423/fnote",
                "css_class": "i-bi:link-45deg",
                "is_link": true
            }
        ]
    },
    create_time: Date.now(),
    update_time: Date.now()
});
// 首页展示文章数量配置
db.getCollection("configs").insertOne({
    _id: "front-post-count",
    typ: "front-post-count",
    "props": {
        "count": 6
    },
    create_time: Date.now(),
    update_time: Date.now()
});
// 支付二维码配置
db.getCollection("configs").insertOne({
    _id: "pay",
    typ: "pay",
    "props": {
        "list": [
            {
                "name": "支付宝收款码",
                "image": "https://blog-1302954944.cos.ap-guangzhou.myqcloud.com/1695022533043.jpg"
            },
            {
                "name": "微信收款码",
                "image": "https://blog-1302954944.cos.ap-guangzhou.myqcloud.com/mm_facetoface_collect_qrcode_1702566807960.png"
            }
        ]
    },
    create_time: Date.now(),
    update_time: Date.now()
});

// ----------------------------
// Collection structure for friends
// ----------------------------
db.getCollection("friends").drop();
db.createCollection("friends");
db.getCollection("friends").createIndex({
    url: NumberInt("1")
}, {
    name: "url_1",
    unique: true
});
// 创建 create_time 降序索引
db.friends.createIndex({ "create_time": -1 });
db.friends.createIndex({ "status": -1 });

// ----------------------------
// Collection structure for message_template
// ----------------------------
db.getCollection("message_template").drop();
db.createCollection("message_template");
// 创建 name 升序索引
db.message_template.createIndex({ "name": 1 });

// ----------------------------
// Documents of message_template
// ----------------------------
db.getCollection("message_template").insertOne({
    _id: "friend",
    name: "friend",
    title: "友链申请通知",
    content: "您好，您的网站有了新的友链申请，详情可前往后台查看。",
    create_time: Date.now(),
    update_time: Date.now(),
    active: 1,
    recipient_type: 0
});
db.getCollection("message_template").insertOne({
    _id: "comment",
    name: "comment",
    title: "文章评论通知",
    content: "您好，您在文章有新的评论，详情请前往后台进行查看。",
    create_time: Date.now(),
    update_time: Date.now(),
    recipient_type: 0,
    active: 1
});

// ----------------------------
// Collection structure for posts
// ----------------------------
db.getCollection("posts").drop();
db.createCollection("posts");
// 创建 create_time 降序索引
db.posts.createIndex({ "create_time": -1 });
// 创建 create_time 升序索引
db.posts.createIndex({ "create_time": 1 });
// 创建 category 单字段索引
db.posts.createIndex({ "categories": 1 });
// 创建 category 和 tags 复合索引
db.posts.createIndex({ "tags": 1 });
// 创建 title 文本索引
db.posts.createIndex({ "title": "text" });

// ----------------------------
// Collection structure for visit_logs
// ----------------------------
db.getCollection("visit_logs").drop();
db.createCollection("visit_logs");
// 创建 create_time 降序索引
db.visit_logs.createIndex({ "create_time": -1 });
EOF
