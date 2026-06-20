// Navicat：先执行 query-missing-file-used-in.js，确认结果后再执行本脚本。
// 本脚本只补充缺失的 post / post-draft / config 引用，不删除任何已有 used_in。
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
                                    { value: { $ifNull: ["$props.website_icon", ""] } },
                                    { value: { $ifNull: ["$props.website_owner_avatar", ""] } },
                                    { value: { $ifNull: ["$props.og_image", ""] } },
                                    { value: { $ifNull: ["$props.content", ""] } },
                                    { value: { $ifNull: ["$props.introduction", ""] } },
                                ],
                                {
                                    $map: {
                                        input: { $ifNull: ["$props.website_records", []] },
                                        as: "value",
                                        in: { value: "$$value" },
                                    },
                                },
                                {
                                    $map: {
                                        input: { $ifNull: ["$props.list", []] },
                                        as: "item",
                                        in: { value: { $ifNull: ["$$item.image", ""] } },
                                    },
                                },
                                {
                                    $map: {
                                        input: { $ifNull: ["$props.list", []] },
                                        as: "item",
                                        in: { value: { $ifNull: ["$$item.cover_img", ""] } },
                                    },
                                },
                                {
                                    $map: {
                                        input: { $ifNull: ["$props.social_info_list", []] },
                                        as: "item",
                                        in: { value: { $ifNull: ["$$item.social_value", ""] } },
                                    },
                                },
                            ],
                        },
                    },
                },
                {
                    $match: {
                        $expr: {
                            $anyElementTrue: {
                                $map: {
                                    input: "$config_file_fields",
                                    as: "field",
                                    in: {
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
                },
                {
                    $project: {
                        _id: 0,
                        entity_id: { $toString: "$_id" },
                        entity_type: { $literal: "config" },
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
            missing_used_in: {
                $map: {
                    input: {
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
                    as: "reference",
                    in: {
                        entity_id: "$$reference.entity_id",
                        entity_type: "$$reference.entity_type",
                    },
                },
            },
        },
    },
    {
        $match: {
            $expr: { $gt: [{ $size: "$missing_used_in" }, 0] },
        },
    },
    {
        $set: {
            used_in: {
                $concatArrays: [
                    { $ifNull: ["$used_in", []] },
                    "$missing_used_in",
                ],
            },
            updated_at: "$$NOW",
        },
    },
    {
        $unset: [
            "post_references",
            "draft_references",
            "config_references",
            "existing_usage_keys",
            "all_references",
            "missing_used_in",
        ],
    },
    {
        $merge: {
            into: "file_meta",
            on: "_id",
            whenMatched: "merge",
            whenNotMatched: "discard",
        },
    },
], { allowDiskUse: true });
