import instance from '@/utils/axios'
import type { LoginRequest } from '@/interfaces/User'

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

export interface WebsiteConfigRequest {
  website_name: string
  live_time: number
  icon: string
  owner_name: string
  owner_profile: string
  owner_picture: string
}

export interface SeoConfigRequest {
  title: string
  description: string
  og_title: string
  og_image: string
  baidu_site_verification: string
  keywords: string
  author: string
  robots: string
}

export const GetWebSite = () => {
  return instance({
    url: '/configs/website',
    method: 'get'
  })
}

export const UpdateWebSite = (req: WebsiteConfigRequest) => {
  return instance({
    url: '/configs/website',
    method: 'put',
    data: req
  })
}

export const AddRecord = (record: string) => {
  return instance({
    url: '/configs/website/records',
    method: 'post',
    data: {
      record: record
    }
  })
}

export const DeleteRecord = (record: string) => {
  return instance({
    url: '/configs/website/records?record=' + record,
    method: 'delete'
  })
}

export const GetSeo = () => {
  return instance({
    url: '/configs/seo',
    method: 'get'
  })
}

export const UpdateSeo = (req: SeoConfigRequest) => {
  return instance({
    url: '/configs/seo',
    method: 'put',
    data: req
  })
}

export const GetComment = () => {
  return instance({
    url: '/configs/comment',
    method: 'get'
  })
}

export interface CommentConfigRequest {
  enable_comment: boolean
}

export const UpdateComment = (req: CommentConfigRequest) => {
  return instance({
    url: '/configs/comment',
    method: 'put',
    data: req
  })
}

export const GetFriend = () => {
  return instance({
    url: '/configs/friend',
    method: 'get'
  })
}

export interface FriendConfigRequest {
  enable_friend_commit: boolean
}

export const UpdateFriend = (req: FriendConfigRequest) => {
  return instance({
    url: '/configs/friend',
    method: 'put',
    data: req
  })
}

export const GetEmail = () => {
  return instance({
    url: '/configs/email',
    method: 'get'
  })
}

export interface EmailConfigRequest {
  host: string
  port: number
  username: string
  password: string
  email: string
}

export const UpdateEmail = (req: EmailConfigRequest) => {
  return instance({
    url: '/configs/email',
    method: 'put',
    data: req
  })
}

export const GetNotice = () => {
  return instance({
    url: '/configs/notice',
    method: 'get'
  })
}

export interface NoticeConfigRequest {
  title: string
  content: string
}

export const UpdateNotice = (req: NoticeConfigRequest) => {
  return instance({
    url: '/configs/notice',
    method: 'put',
    data: req
  })
}

export const UpdateNoticeEnabled = (enabled: boolean) => {
  return instance({
    url: '/configs/notice/enabled',
    method: 'put',
    data: {
      enabled: enabled
    }
  })
}

export const GetFrontPostCount = () => {
  return instance({
    url: '/configs/front-post-count',
    method: 'get'
  })
}

export interface FrontPostCountConfigRequest {
  count: number
}

export const UpdateFrontPostCount = (req: FrontPostCountConfigRequest) => {
  return instance({
    url: '/configs/front-post-count',
    method: 'put',
    data: req
  })
}

export const GetPay = () => {
  return instance({
    url: '/configs/pay',
    method: 'get'
  })
}

export const AddPay = (req: PayConfigRequest) => {
  return instance({
    url: '/configs/pay',
    method: 'post',
    data: req
  })
}

export const DeletePay = (name: string, image: string) => {
  return instance({
    url: `/configs/pay/${name}?image=${image}`,
    method: 'delete'
  })
}

export const GetSocial = () => {
  return instance({
    url: '/configs/social',
    method: 'get'
  })
}

export interface SocialConfigRequest {
  social_name: string
  social_value: string
  css_class: string
  is_link: boolean
}

export const AddSocial = (req: SocialConfigRequest) => {
  return instance({
    url: '/configs/social',
    method: 'post',
    data: req
  })
}

export const DeleteSocial = (id: string) => {
  return instance({
    url: `/configs/social/${id}`,
    method: 'delete'
  })
}

export interface SocialConfigRequest {}

export const UpdateSocial = (id: string, req: SocialConfigRequest) => {
  return instance({
    url: `/configs/social/${id}`,
    method: 'put',
    data: req
  })
}

export interface InitReq {
  website_name: string
  website_icon: string
  website_owner: string
  website_owner_profile: string
  website_owner_avatar: string
  website_domain: string
  website_owner_email: string
  email_server: {
    host: string
    port: string
    username: string
    password: string
    email: string
  }
}

export const isInit = () => {
  return instance({
    url: '/check-initialization',
    method: 'get',
  })
}
