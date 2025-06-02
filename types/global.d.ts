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

type TaskModel = {
  id: string;
  title: string;
  description?: string;
  creator_id: number;
  deadline: string;
  completed: boolean;
  created_at: string;
};
type Task = TaskModel;
type CreateTaskResponseData = Task;
type GetTaskByIdResponseData = Task;
