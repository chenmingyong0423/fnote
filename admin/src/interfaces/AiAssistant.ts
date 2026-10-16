import instance from '@/utils/axios'
import type { IBaseResponse, IResponse } from '@/interfaces/Common'

export type AIAPIType = 'openai_responses' | 'openai_chat_completions'

export interface AIProvider {
  id: string
  name: string
  api_type: AIAPIType
  base_url: string
  model: string
  has_api_key: boolean
  max_output_tokens?: number | null
  reasoning_effort?: string
  thinking_mode?: '' | 'enabled' | 'disabled'
}

export interface AIAssistantConfig {
  enabled: boolean
  assistant_name: string
  system_prompt: string
  enabled_tools: string[]
  active_provider_id: string
  providers: AIProvider[]
}

export type AIProviderRequest = Omit<AIProvider, 'has_api_key'> & { api_key?: string }

export type UpdateAIAssistantConfigRequest = Omit<AIAssistantConfig, 'providers'> & {
  providers: AIProviderRequest[]
}

export interface AIProviderTestResult {
  success: boolean
  message: string
  reply?: string
  duration_ms: number
}

export const GetAIAssistantConfig = () =>
  instance.get<IResponse<AIAssistantConfig>>('/configs/ai-assistant')

export const UpdateAIAssistantConfig = (data: UpdateAIAssistantConfigRequest) =>
  instance.put<IBaseResponse>('/configs/ai-assistant', data)

export const TestAIProvider = (data: AIProviderRequest, signal?: AbortSignal) =>
  instance.post<IResponse<AIProviderTestResult>>('/configs/ai-assistant/test-provider', data, {
    signal,
    timeout: 35000
  })
