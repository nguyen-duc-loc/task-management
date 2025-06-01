import { Metadata } from "next";
import React from "react";
import SigninForm from "./_components/SigninForm";

export const metadata: Metadata = {
  title: "Sign in",
  description: "Sign in to task management system",
};

const SigninPage = () => {
  return <SigninForm />;
};

export default SigninPage;
