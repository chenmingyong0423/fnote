export interface Comment{
  id: string;
  post_info: {
    post_id: string;
    post_title: string;
  }
  content: string;
  user_info: {
    username: string;
    email: string;
    website?: string;
    ip: string;
  }
  fid?: string;
  type: number;
}