const ROUTES = {
  dashboard: "/",
  signin: "/signin",
  signup: "/signup",
  task: (id: string) => `/task/${id}`,
};

export default ROUTES;
