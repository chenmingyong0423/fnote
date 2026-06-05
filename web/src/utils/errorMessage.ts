type ErrorLike = {
  message?: string;
  status?: number;
  statusCode?: number;
};

const DEFAULT_ERROR_MESSAGE = "操作失败，请稍后重试";

const STATUS_MESSAGE_MAP: Record<number, string> = {
  400: "请求参数有误",
  401: "请先登录后再操作",
  403: "当前无权限执行该操作",
  404: "请求的内容不存在",
  429: "请求过于频繁，请稍后再试",
  500: "服务异常，请稍后重试",
  502: "网关异常，请稍后重试",
  503: "服务暂不可用，请稍后重试",
  504: "服务响应超时，请稍后重试",
};

const EXACT_MESSAGE_MAP: Record<string, string> = {
  "ip is empty": "IP 不能为空",
  "website format is invalid": "网站地址格式不正确，请使用 https:// 开头",
  "comment module is closed": "评论功能暂未开启",
  "post not found": "文章不存在",
  "comments are disabled for this post": "该文章已关闭评论",
  "email format is incorrect": "邮箱格式不正确",
  "friend module is close": "友链功能暂未开启",
  "already applied for friendship, please wait for review": "您已申请过友链，请等待审核",
};

const KEYWORD_RULES: Array<{ keywords: string[]; message: string }> = [
  { keywords: ["ip", "empty"], message: "IP 不能为空" },
  { keywords: ["website", "invalid"], message: "网站地址格式不正确，请使用 https:// 开头" },
  { keywords: ["comment", "closed"], message: "评论功能暂未开启" },
  { keywords: ["post", "not", "found"], message: "文章不存在" },
  { keywords: ["comments", "disabled"], message: "该文章已关闭评论" },
  { keywords: ["email", "incorrect"], message: "邮箱格式不正确" },
  { keywords: ["friend", "module", "close"], message: "友链功能暂未开启" },
  { keywords: ["applied", "friendship", "review"], message: "您已申请过友链，请等待审核" },
];

function extractError(error: unknown): ErrorLike {
  if (error && typeof error === "object") {
    return error as ErrorLike;
  }
  return {};
}

function normalizeMessage(value: string | undefined): string {
  return (value ?? "")
    .trim()
    .toLowerCase()
    .replace(/[.!?]+$/g, "");
}

export function getUserFriendlyError(error: unknown): string {
  const normalized = extractError(error);
  const backendMessage = normalizeMessage(normalized.message);

  if (backendMessage) {
    const exactMatched = EXACT_MESSAGE_MAP[backendMessage];
    if (exactMatched) {
      return exactMatched;
    }

    const keywordMatched = KEYWORD_RULES.find((rule) =>
      rule.keywords.every((keyword) => backendMessage.includes(keyword))
    );
    if (keywordMatched) {
      return keywordMatched.message;
    }
  }

  const status = normalized.status ?? normalized.statusCode;
  if (typeof status === "number" && STATUS_MESSAGE_MAP[status]) {
    return STATUS_MESSAGE_MAP[status];
  }

  return DEFAULT_ERROR_MESSAGE;
}
