import instance from '@/utils/axios'
import type { VNode } from 'vue'

export interface AssetFolderVO {
  id: string
  name: string
  // icon: () => h(PictureOutlined),
  icon: VNode
  support_delete: boolean
  support_edit: boolean
  support_add: boolean
  children: AssetFolderVO[]
}

export interface AddAssetFolderRequest {
  id: string
  name: string
  asset_type: string
  type: string
  support_delete: boolean
  support_edit: boolean
  support_add: boolean
}

export interface AssetVO {
  id: string
  title: string
  content: string
  description: string
  asset_type: string
  type: string
  metadata: any
}

export interface AssetRequest {
  id: string
  title: string
  content: string
  description: string
  asset_type: string
  type: string
  metadata: any
}

export const GetAssetFolderList = (assetType: string, typ: string) => {
  return instance({
    url: '/assets/folders',
    method: 'get',
    params: {
      assetType: assetType,
      type: typ
    }
  })
}

export const AddAssetFolder = (assetFolderRequest: AddAssetFolderRequest) => {
  return instance({
    url: '/assets/folders',
    method: 'post',
    data: assetFolderRequest
  })
}

export const EditAssetFolder = (id: string, name: string) => {
  return instance({
    url: `/assets/folders/${id}`,
    method: 'put',
    data: {
      name: name
    }
  })
}
export const EditAssetFolderName = (id: string, name: string) => {
  return instance({
    url: `/assets/folders/${id}/name`,
    method: 'put',
    data: {
      name: name
    }
  })
}

export const DeleteAssetFolder = (id: string) => {
  return instance({
    url: `/assets/folders/${id}`,
    method: 'delete'
  })
}

export const GetAssetList = (folderId: string) => {
  return instance({
    url: `/assets/folders/${folderId}/assets`,
    method: 'get'
  })
}

export const AddAsset = (folderId: string, assetRequest: AssetRequest) => {
  return instance({
    url: `/assets/folders/${folderId}/assets`,
    method: 'post',
    data: assetRequest
  })
}
