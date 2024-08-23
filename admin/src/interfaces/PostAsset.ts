import instance from '@/utils/axios'

export interface AssetFolderVO {
  id: string
  name: string
  // icon: () => h(PictureOutlined),
  icon: string
  support_upload: boolean
  support_delete: boolean
  support_edit: boolean
}

export const GetAssetFolderList = (assetType:string, typ:string) => {
  return instance({
    url: '/assets/folders',
    method: 'get',
    params: {
      "assetType": assetType,
      "type": typ
    }
  })
}