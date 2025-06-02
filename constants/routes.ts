const ROUTES = {
  dashboard: "/",
  signin: "/signin",
  signup: "/signup",
  task: (id: string) => `/task/${id}`,
  editTask: (id: string) => `/task/${id}/edit`,
};

export default ROUTES;
