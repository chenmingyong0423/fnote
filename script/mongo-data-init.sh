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
db.getCollection("configs").createIndex({
    typ: NumberInt("1")
}, {
    name: "unique_typ",
    unique: true
});
// 站点信息
db.getCollection("configs").insertOne({
  "created_at": new Date(),
  "props": {
    // 站点名称
    "website_name": "",
    // 站点图标
    "website_icon": "",
    // 站长昵称
    "website_owner": "",
    // 站长简介
    "website_owner_profile": "",
    // 站长头像
    "website_owner_avatar": "",
    // 站点运行时间
    "website_runtime": new Date(),
    // 站点备案号
    "website_records": [],
    // 是否完成初始化
    "website_init": false
  },
  "typ": "website",
  "updated_at": new Date()
});
// seo 配置
db.getCollection("configs").insertOne({
    "typ": "seo meta",
    "props": {
        "title": "",
        "og_title": "",
        "description": "",
        "og_image": "",
        "keywords": "",
        "author": "",
        "robots": "",
        "third_party_site_verification": []
    },
    created_at: new Date(),
    updated_at: new Date()
});
// 文章索引推送配置
db.getCollection("configs").insertOne({
    "typ": "post index",
    "props": {},
    created_at: new Date(),
    updated_at: new Date()
});
// 第三方站点验证配置
db.getCollection("configs").insertOne({
    "typ": "third party site verification",
    "props": {
        "list": []
    },
    created_at: new Date(),
    updated_at: new Date()
});

// 评论开关配置
db.getCollection("configs").insertOne({
    typ: "comment",
    props: {
        enable_comment: true
    },
    created_at: new Date(),
    updated_at: new Date()
});
// 友链开关配置
db.getCollection("configs").insertOne({
    typ: "friend",
    props: {
        enable_friend_commit: false,
        introduction: ''
    },
    created_at: new Date(),
    updated_at: new Date()
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
    created_at: new Date(),
    updated_at: new Date()
});
// 公告配置
db.getCollection("configs").insertOne({
    typ: "notice",
    "props": {
        "title": "暂无最新公告",
        "content": "暂无最新公告",
        "publish_time": new Date(),
        "enabled": true
    },
    created_at: new Date(),
    updated_at: new Date()
});
// 社交外链配置
db.getCollection("configs").insertOne({
    typ: "social",
    "props": {
        "social_info_list": []
    },
    created_at: new Date(),
    updated_at: new Date()
});
// 首页展示文章数量配置
db.getCollection("configs").insertOne({
    typ: "front-post-count",
    "props": {
        "count": 6
    },
    created_at: new Date(),
    updated_at: new Date()
});
// 支付二维码配置
db.getCollection("configs").insertOne({
    typ: "pay",
    "props": {
        "list": []
    },
    created_at: new Date(),
    updated_at: new Date()
});

// 管理员配置
db.getCollection("configs").insertOne({
    typ: "admin",
    "props": {
        "username": "",
        "password": ""
    },
    created_at: new Date(),
    updated_at: new Date()
});

// 轮播图数据
db.getCollection("configs").insertOne({
    typ: "carousel",
    "props": {
        "list": []
    },
    created_at: new Date(),
    updated_at: new Date()
});

// friends
db.createCollection("friends");
db.getCollection("friends").createIndex({ status: 1, created_at: 1 })
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
    created_at: new Date(),
    updated_at: new Date(),
    recipient_type: 0,
    active: 1
});

db.getCollection("message_templates").insertOne({
    name: "user-comment-approval",
    title: "评论审核通过通知",
    content: "您好，您在 %s 文章中发表的评论已通过审核。",
    created_at: new Date(),
    updated_at: new Date(),
    recipient_type: 1,
    active: 1
});

db.getCollection("message_templates").insertOne({
    name: "user-comment-disapproval",
    title: "评论被驳回通知",
    content: "您好，您在 %s 文章中发表的评论未通过审核，原因：%s",
    created_at: new Date(),
    updated_at: new Date(),
    recipient_type: 1,
    active: 1
});

db.getCollection("message_templates").insertOne({
    name: "user-comment-reply",
    title: "评论被回复通知",
    content: "您好，您在 %s 文章中发表的评论有新的回复。",
    created_at: new Date(),
    updated_at: new Date(),
    recipient_type: 1,
    active: 1
});

db.getCollection("message_templates").insertOne({
    name: "friend",
    title: "友链申请通知",
    content: "您好，您的网站有了新的友链申请，详情可前往后台查看。",
    created_at: new Date(),
    updated_at: new Date(),
    active: 1,
    recipient_type: 0
});

db.getCollection("message_templates").insertOne({
    name: "friend-approval",
    title: "友链申请通过通知",
    content: "您好，您在 %s 网站里提交的友链申请已通过审核并展示在页面上。",
    created_at: new Date(),
    updated_at: new Date(),
    recipient_type: 1,
    active: 1
});

db.getCollection("message_templates").insertOne({
    name: "friend-rejection",
    title: "友链申请不通过通知",
    content: "您好，您在 %s 网站里提交的友链申请未通过审核，原因：%s",
    created_at: new Date(),
    updated_at: new Date(),
    recipient_type: 1,
    active: 1
});

// posts
db.createCollection("posts");
// 创建 created_at 降序索引
db.getCollection("posts").createIndex({ "created_at": -1 });
// 创建 created_at 升序索引
db.getCollection("posts").createIndex({ "created_at": 1 });
// 创建 category 单字段索引
db.getCollection("posts").createIndex({ "categories": 1 });
// 创建 tags 单字段索引
db.getCollection("posts").createIndex({ "tags": 1 });
// 创建 title 文本索引
db.getCollection("posts").createIndex({ "title": "text" });

// visit_logs
db.createCollection("visit_logs");
// 创建 created_at 降序索引
db.getCollection("visit_logs").createIndex({ "created_at": -1 });

// file_meta
db.createCollection("file_meta");
// 为 file_name  创建唯一索引
db.getCollection("file_meta").createIndex({
    file_name: NumberInt("1")
}, {
    name: "unique_file_name",
    unique: true
});

// count_stats
db.createCollection("count_stats")
db.getCollection("count_stats").createIndex({
    type: NumberInt("1"),
}, {
    name: "unique_type",
    unique: true
});
db.getCollection("count_stats").insertMany([
    {
        "type": "PostCount",
        "count": 0,
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        "type": "CategoryCount",
        "count": 0,
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        "type": "TagCount",
        "count": 0,
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        "type": "CommentCount",
        "count": 0,
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        "type": "LikeCount",
        "count": 0,
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        "type": "WebsiteViewCount",
        "count": 0,
        created_at: new Date(),
        updated_at: new Date()
    }
])
// post-likes
db.createCollection("post_likes");
db.post_likes.createIndex({ "post_id": 1, "ip": 1 }, { "unique": true })
EOF
