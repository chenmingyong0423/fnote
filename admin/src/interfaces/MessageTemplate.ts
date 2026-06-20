import instance from '@/utils/axios'

export interface MessageTemplate {
  id: string
  type: string
  name: string
  title: string
  content: string
  active: boolean
  is_default: boolean
  recipient_type: number
  variables: string[]
  created_at: number
  updated_at: number
}

export interface UpdateMessageTemplateRequest {
  name: string
  title: string
  content: string
}

export interface CreateMessageTemplateRequest extends UpdateMessageTemplateRequest {
  type: string
  active: boolean
}

export const GetMessageTemplates = () => {
  return instance({
    url: '/message-templates',
    method: 'get'
  })
}

export const CreateMessageTemplate = (request: CreateMessageTemplateRequest) => {
  return instance({ url: '/message-templates', method: 'post', data: request })
}

export const UpdateMessageTemplate = (id: string, request: UpdateMessageTemplateRequest) => {
  return instance({
    url: `/message-templates/${id}`,
    method: 'put',
    data: request
  })
}

export const UpdateMessageTemplateActive = (id: string, active: boolean) => {
  return instance({
    url: `/message-templates/${id}/active`,
    method: 'put',
    data: { active }
  })
}

export const SetDefaultMessageTemplate = (id: string) => {
  return instance({ url: `/message-templates/${id}/default`, method: 'put' })
}

export const DeleteMessageTemplate = (id: string) => {
  return instance({ url: `/message-templates/${id}`, method: 'delete' })
}
