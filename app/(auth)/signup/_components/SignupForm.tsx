"use client";

import { signup } from "@/api/actions/auth";
import AuthForm from "@/components/forms/AuthForm";
import { SignupSchema, SignupData } from "@/validation";
import React from "react";

const SignupForm = () => {
  return (
    <AuthForm
      schema={SignupSchema}
      defaultValues={{
        username: "",
        password: "",
        confirmPassword: "",
      }}
      onSubmit={(data: SignupData) => signup(data)}
      formType="SIGN_UP"
    />
  );
};

export default SignupForm;
