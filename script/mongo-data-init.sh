mongosh -- "$MONGO_INITDB_DATABASE" <<EOF
db = db.getSiblingDB('$MONGO_INITDB_DATABASE')
db.auth('$MONGO_USERNAME', '$MONGO_PASSWORD');

// categories
db.createCollection("categories");
db.getCollection("categories").createIndex({
    name: NumberInt("1")
}, {
    name: "unique_name",
    unique: true
});
db.getCollection("categories").createIndex({
    route: NumberInt("1")
}, {
    name: "unique_route",
    unique: true
});

// tags
db.createCollection("tags");
db.getCollection("tags").createIndex({
    name: NumberInt("1")
}, {
    name: "unique_name",
    unique: true
});
db.getCollection("tags").createIndex({
    route: NumberInt("1")
}, {
    name: "unique_route",
    unique: true
});

// comments
db.createCollection("comments");
db.getCollection("comments").createIndex({ "post_info.post_id": 1 });
db.getCollection("comments").createIndex({ "status": -1 });

// configs
db.createCollection("configs");
db.getCollection("configs").createIndex({ "typ": 1 });
// 站点信息
db.getCollection("configs").insertOne({
  "create_time": Math.floor(new Date().getTime() / 1000),
  "props": {
    "name": "fnote",
    "post_count": 0,
    "category_count": 0,
    "view_count": 0,
    "live_time": Math.floor(new Date().getTime() / 1000),
    "icon": "",
    "domain": "",
    "records": []
  },
  "typ": "website",
  "update_time": Math.floor(new Date().getTime() / 1000)
});
// 站长信息
db.getCollection("configs").insertOne(
  {
    "create_time": Math.floor(new Date().getTime() / 1000),
    "props": {
      "name": "fnote user",
      "profile": "请及时前往后台修改站点和站长等相关配置，以便正常使用。",
      "picture": ""
    },
    "typ": "owner",
    "update_time": Math.floor(new Date().getTime() / 1000)
  }
)
// seo 配置
db.getCollection("configs").insertOne({
    "typ": "seo meta",
    "props": {
        "title": "fnote",
        "og_title": "fnote",
        "description": "fnote",
        "og_image": "",
        "baidu_site_verification": "",
        "keywords": "fnote,blog,BLOG",
        "author": "fnote",
        "robots": "fnote,blog"
    },
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000)
});
// 评论开关配置
db.getCollection("configs").insertOne({
    typ: "comment",
    props: {
        enable_comment: true
    },
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000)
});
// 友链开关配置
db.getCollection("configs").insertOne({
    typ: "friend",
    props: {
        enable_friend_commit: false
    },
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000)
});
// 邮件配置
db.getCollection("configs").insertOne({
    "typ": "email",
    "props": {
        "host": "",
        "port": 0,
        "username": "",
        "password": "",
        "email": ""
    },
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000)
});
// 公告配置
db.getCollection("configs").insertOne({
    typ: "notice",
    "props": {
        "title": "暂无最新公告",
        "content": "暂无最新公告",
        "publish_time": Math.floor(new Date().getTime() / 1000),
        "enabled": true
    },
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000)
});
// 社交外链配置
db.getCollection("configs").insertOne({
    typ: "social",
    "props": {
        "social_info_list": []
    },
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000)
});
// 首页展示文章数量配置
db.getCollection("configs").insertOne({
    typ: "front-post-count",
    "props": {
        "count": 6
    },
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000)
});
// 支付二维码配置
db.getCollection("configs").insertOne({
    typ: "pay",
    "props": {
        "list": []
    },
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000)
});

// friends
db.createCollection("friends");
db.getCollection("friends").createIndex({ status: 1, create_time: 1 })
db.getCollection("friends").createIndex({
    url: NumberInt("1")
}, {
    name: "unique_url",
    unique: true
});

// message_templates
db.createCollection("message_templates");
// 创建 name 升序索引
db.getCollection("message_templates").createIndex({
    name: NumberInt("1")
}, {
    name: "unique_name",
    unique: true
});


db.getCollection("message_templates").insertOne({
    name: "comment",
    title: "文章评论通知",
    content: "您好，您在文章有新的评论，详情请前往后台进行查看。",
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000),
    recipient_type: 0,
    active: 1
});

db.getCollection("message_templates").insertOne({
    name: "user-comment-approval",
    title: "评论审核通过通知",
    content: "您好，您在 %s 文章中发表的评论已通过审核。",
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000),
    recipient_type: 1,
    active: 1
});

db.getCollection("message_templates").insertOne({
    name: "user-comment-disapproval",
    title: "评论被驳回通知",
    content: "您好，您在 %s 文章中发表的评论未通过审核，原因：%s",
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000),
    recipient_type: 1,
    active: 1
});

db.getCollection("message_templates").insertOne({
    name: "user-comment-reply",
    title: "评论被回复通知",
    content: "您好，您在 %s 文章中发表的评论有新的回复。",
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000),
    recipient_type: 1,
    active: 1
});

db.getCollection("message_templates").insertOne({
    name: "friend",
    title: "友链申请通知",
    content: "您好，您的网站有了新的友链申请，详情可前往后台查看。",
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000),
    active: 1,
    recipient_type: 0
});

db.getCollection("message_templates").insertOne({
    name: "friend-approval",
    title: "友链申请通过通知",
    content: "您好，您在 %s 网站里提交的友链申请已通过审核并展示在页面上。",
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000),
    recipient_type: 1,
    active: 1
});

db.getCollection("message_templates").insertOne({
    name: "friend-rejection",
    title: "友链申请不通过通知",
    content: "您好，您在 %s 网站里提交的友链申请未通过审核，原因：%s",
    create_time: Math.floor(new Date().getTime() / 1000),
    update_time: Math.floor(new Date().getTime() / 1000),
    recipient_type: 1,
    active: 1
});

// posts
db.createCollection("posts");
// 创建 create_time 降序索引
db.getCollection("posts").createIndex({ "create_time": -1 });
// 创建 create_time 升序索引
db.getCollection("posts").createIndex({ "create_time": 1 });
// 创建 category 单字段索引
db.getCollection("posts").createIndex({ "categories": 1 });
// 创建 tags 单字段索引
db.getCollection("posts").createIndex({ "tags": 1 });
// 创建 title 文本索引
db.getCollection("posts").createIndex({ "title": "text" });

// visit_logs
db.getCollection("visit_logs").drop();
db.createCollection("visit_logs");
// 创建 create_time 降序索引
db.getCollection("visit_logs").createIndex({ "create_time": -1 });
EOF
