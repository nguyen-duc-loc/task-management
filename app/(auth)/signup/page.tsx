import { Metadata } from "next";
import React from "react";
import SignupForm from "./_components/SignupForm";

export const metadata: Metadata = {
  title: "Sign up",
  description: "Sign up to task management system",
};

const SignupPage = () => {
  return <SignupForm />;
};

export default SignupPage;
