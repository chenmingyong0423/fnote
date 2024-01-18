export interface WebsiteConfig {
  website_name: string
  icon: string
  live_time: number
  records: string[]
  owner_name: string
  owner_profile: string
  owner_picture: string
}

export interface SeoConfig {
  title: string
  description: string
  og_title: string
  og_image: string
  baidu_site_verification: string
  keywords: string
  author: string
  robots: string
}

export interface CommentConfig {
  enable_comment: boolean
}

export interface FriendConfig {
  enable_friend_commit: boolean
}

export interface EmailConfig {
  host: string
  port: number
  username: string
  password: string
  email: string
}

export interface NoticeConfig {
  title: string
  content: string
  enabled: boolean
  publish_time: number
}

export interface FrontPostCountConfig {
  count: number
}

export interface PayConfig {
  name: string
  image: string
}

export interface PayConfigRequest {
  name: string
  image: string
}

export interface SocialConfig {
  id: string
  social_name: string
  social_value: string
  css_class: string
  is_link: boolean
}

export interface SocialConfigReq {
  social_name: string
  social_value: string
  css_class: string
  is_link: boolean
}
