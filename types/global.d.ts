type ActionResponse<T = null> = {
  success: boolean;
  data?: T;
  error?: string;
};

type UserModel = {
  id: string;
  username: string;
  hashed_password: string;
  created_at: string;
};
type User = Omit<UserModel, "hashed_password">;
type SignupResponseData = User;
type SigninResponseData = {
  access_token: string;
  access_token_expire_at: string;
  user: User;
};
