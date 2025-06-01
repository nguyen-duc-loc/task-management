"use client";

import { signin } from "@/api/actions/auth";
import AuthForm from "@/components/forms/AuthForm";
import { SigninSchema, SigninData } from "@/validation";
import React from "react";

const SigninForm = () => {
  return (
    <AuthForm
      schema={SigninSchema}
      defaultValues={{
        username: "",
        password: "",
      }}
      onSubmit={(data: SigninData) => signin(data)}
      formType="SIGN_IN"
    />
  );
};

export default SigninForm;
