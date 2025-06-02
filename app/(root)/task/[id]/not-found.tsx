import React from "react";

import ErrorPage from "@/components/error-page";

const NotFoundTask = () => {
  return (
    <ErrorPage
      errorCode={404}
      message="Something went wrong"
      description="Sorry we were unable to find the task you are looking for"
    />
  );
};

export default NotFoundTask;
