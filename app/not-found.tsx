import React from "react";

import ErrorPage from "@/components/error-page";

const NotFound = () => {
  return (
    <div className="flex h-screen items-center justify-center">
      <ErrorPage
        errorCode={404}
        message="Something went wrong"
        description="Sorry we were unable to find that page"
      />
    </div>
  );
};

export default NotFound;
