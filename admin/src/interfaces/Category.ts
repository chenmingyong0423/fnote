export interface ICategory {
  id: string;
  name: string;
  route: string;
  description: string;
  disabled: boolean;
  create_time: number;
  update_time: number;
}

export interface CategoryRequest {
  name: string;
  route: string;
  description: string;
  disabled: boolean;
}