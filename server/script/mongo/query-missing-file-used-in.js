// Navicat：选择 fnote 数据库，新建查询，粘贴整个脚本后执行。
// 本脚本只查询，不修改数据。
db.getCollection("file_meta").aggregate([
    {
        $match: {
            url: { $type: "string", $ne: "" },
        },
    },
    {
        $lookup: {
            from: "posts",
            let: { fileUrl: "$url" },
            pipeline: [
                {
                    $match: {
                        $expr: {
                            $or: [
                                { $gte: [{ $indexOfCP: [{ $ifNull: ["$content", ""] }, "$$fileUrl"] }, 0] },
                                { $gte: [{ $indexOfCP: [{ $ifNull: ["$cover_img", ""] }, "$$fileUrl"] }, 0] },
                            ],
                        },
                    },
                },
                {
                    $project: {
                        _id: 0,
                        entity_id: { $toString: "$_id" },
                        entity_type: { $literal: "post" },
                        title: { $ifNull: ["$title", ""] },
                        locations: {
                            $concatArrays: [
                                {
                                    $cond: [
                                        { $gte: [{ $indexOfCP: [{ $ifNull: ["$content", ""] }, "$$fileUrl"] }, 0] },
                                        ["content"],
                                        [],
                                    ],
                                },
                                {
                                    $cond: [
                                        { $gte: [{ $indexOfCP: [{ $ifNull: ["$cover_img", ""] }, "$$fileUrl"] }, 0] },
                                        ["cover_img"],
                                        [],
                                    ],
                                },
                            ],
                        },
                    },
                },
            ],
            as: "post_references",
        },
    },
    {
        $lookup: {
            from: "post_draft",
            let: { fileUrl: "$url" },
            pipeline: [
                {
                    $match: {
                        $expr: {
                            $or: [
                                { $gte: [{ $indexOfCP: [{ $ifNull: ["$content", ""] }, "$$fileUrl"] }, 0] },
                                { $gte: [{ $indexOfCP: [{ $ifNull: ["$cover_img", ""] }, "$$fileUrl"] }, 0] },
                            ],
                        },
                    },
                },
                {
                    $project: {
                        _id: 0,
                        entity_id: { $toString: "$_id" },
                        entity_type: { $literal: "post-draft" },
                        title: { $ifNull: ["$title", ""] },
                        locations: {
                            $concatArrays: [
                                {
                                    $cond: [
                                        { $gte: [{ $indexOfCP: [{ $ifNull: ["$content", ""] }, "$$fileUrl"] }, 0] },
                                        ["content"],
                                        [],
                                    ],
                                },
                                {
                                    $cond: [
                                        { $gte: [{ $indexOfCP: [{ $ifNull: ["$cover_img", ""] }, "$$fileUrl"] }, 0] },
                                        ["cover_img"],
                                        [],
                                    ],
                                },
                            ],
                        },
                    },
                },
            ],
            as: "draft_references",
        },
    },
    {
        $lookup: {
            from: "configs",
            let: { fileUrl: "$url" },
            pipeline: [
                {
                    $set: {
                        config_file_fields: {
                            $concatArrays: [
                                [
                                    { location: "props.website_icon", value: { $ifNull: ["$props.website_icon", ""] } },
                                    { location: "props.website_owner_avatar", value: { $ifNull: ["$props.website_owner_avatar", ""] } },
                                    { location: "props.og_image", value: { $ifNull: ["$props.og_image", ""] } },
                                    { location: "props.content", value: { $ifNull: ["$props.content", ""] } },
                                    { location: "props.introduction", value: { $ifNull: ["$props.introduction", ""] } },
                                ],
                                {
                                    $map: {
                                        input: { $ifNull: ["$props.website_records", []] },
                                        as: "value",
                                        in: { location: "props.website_records", value: "$$value" },
                                    },
                                },
                                {
                                    $map: {
                                        input: { $ifNull: ["$props.list", []] },
                                        as: "item",
                                        in: { location: "props.list.image", value: { $ifNull: ["$$item.image", ""] } },
                                    },
                                },
                                {
                                    $map: {
                                        input: { $ifNull: ["$props.list", []] },
                                        as: "item",
                                        in: { location: "props.list.cover_img", value: { $ifNull: ["$$item.cover_img", ""] } },
                                    },
                                },
                                {
                                    $map: {
                                        input: { $ifNull: ["$props.social_info_list", []] },
                                        as: "item",
                                        in: {
                                            location: "props.social_info_list.social_value",
                                            value: { $ifNull: ["$$item.social_value", ""] },
                                        },
                                    },
                                },
                            ],
                        },
                    },
                },
                {
                    $set: {
                        matched_config_fields: {
                            $filter: {
                                input: "$config_file_fields",
                                as: "field",
                                cond: {
                                    $gte: [
                                        {
                                            $indexOfCP: [
                                                {
                                                    $convert: {
                                                        input: "$$field.value",
                                                        to: "string",
                                                        onError: "",
                                                        onNull: "",
                                                    },
                                                },
                                                "$$fileUrl",
                                            ],
                                        },
                                        0,
                                    ],
                                },
                            },
                        },
                    },
                },
                {
                    $match: {
                        $expr: { $gt: [{ $size: "$matched_config_fields" }, 0] },
                    },
                },
                {
                    $project: {
                        _id: 0,
                        entity_id: { $toString: "$_id" },
                        entity_type: { $literal: "config" },
                        title: { $concat: ["配置：", { $ifNull: ["$typ", "unknown"] }] },
                        locations: {
                            $setUnion: [
                                {
                                    $map: {
                                        input: "$matched_config_fields",
                                        as: "field",
                                        in: "$$field.location",
                                    },
                                },
                                [],
                            ],
                        },
                    },
                },
            ],
            as: "config_references",
        },
    },
    {
        $set: {
            existing_usage_keys: {
                $map: {
                    input: { $ifNull: ["$used_in", []] },
                    as: "usage",
                    in: { $concat: ["$$usage.entity_type", "|", "$$usage.entity_id"] },
                },
            },
            all_references: { $concatArrays: ["$post_references", "$draft_references", "$config_references"] },
        },
    },
    {
        $set: {
            missing_references: {
                $filter: {
                    input: "$all_references",
                    as: "reference",
                    cond: {
                        $not: [
                            {
                                $in: [
                                    { $concat: ["$$reference.entity_type", "|", "$$reference.entity_id"] },
                                    "$existing_usage_keys",
                                ],
                            },
                        ],
                    },
                },
            },
        },
    },
    {
        $match: {
            $expr: { $gt: [{ $size: "$missing_references" }, 0] },
        },
    },
    {
        $project: {
            _id: 1,
            file_id: 1,
            file_name: 1,
            url: 1,
            current_used_in: { $ifNull: ["$used_in", []] },
            missing_references: 1,
        },
    },
], { allowDiskUse: true });
