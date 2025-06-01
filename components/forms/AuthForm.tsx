"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { IconArrowNarrowRight } from "@tabler/icons-react";
import Link from "next/link";
import React from "react";
import {
  DefaultValues,
  Field,
  FieldValues,
  Path,
  useForm,
} from "react-hook-form";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  FormField,
  FormItem,
  FormLabel,
  FormControl,
  FormMessage,
  Form,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import Spinner from "@/components/ui/spinner";
import { ZodType } from "zod";
import ROUTES from "@/constants/routes";

interface AuthFormProps<T extends FieldValues> {
  schema: ZodType<T>;
  defaultValues: T;
  onSubmit: (data: T) => Promise<ActionResponse<unknown>>;
  formType: "SIGN_IN" | "SIGN_UP";
}

const AuthForm = <T extends FieldValues>({
  schema,
  defaultValues,
  onSubmit,
  formType,
}: AuthFormProps<T>) => {
  const form = useForm({
    resolver: zodResolver(schema),
    defaultValues: defaultValues as DefaultValues<T>,
  });

  const handleSubmit = async (data: T) => {
    const response = await onSubmit(data);
    if (!response.success) {
      toast.error(response.error);
    } else {
      form.reset();
    }
  };

  const isSubmitting = form.formState.isSubmitting;

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(handleSubmit)}
        className="max-sm:px-8 w-md"
      >
        <Card>
          <CardHeader>
            <CardTitle className="text-2xl">
              Welcome {formType === "SIGN_IN" ? "back" : "to task management"}
            </CardTitle>
            <CardDescription className="!mt-0">
              Enter your information below to{" "}
              {formType === "SIGN_IN" ? "sign in to" : "create"} your account
            </CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4">
            {Object.keys(defaultValues).map((field) => (
              <FormField
                key={field}
                control={form.control}
                name={field as Path<T>}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="capitalize">
                      {field.name === "confirmPassword"
                        ? "confirm password"
                        : field.name}
                    </FormLabel>
                    <div className="space-y-2">
                      <FormControl>
                        <Input
                          type={field.name === "username" ? "text" : "password"}
                          {...field}
                          disabled={isSubmitting}
                        />
                      </FormControl>
                    </div>
                    <FormMessage />
                  </FormItem>
                )}
              />
            ))}
          </CardContent>
          <CardFooter className="flex-col gap-4">
            <Button className="w-full" type="submit" disabled={isSubmitting}>
              {isSubmitting && <Spinner />}
              {formType === "SIGN_IN" ? "Sign in" : "Sign up"}
            </Button>
            <p className="flex flex-wrap items-center gap-1 text-sm">
              {formType === "SIGN_IN" ? "Don't" : "Already"} have an account?{" "}
              <Link
                href={formType === "SIGN_IN" ? ROUTES.signup : ROUTES.signin}
                className="flex p-0 font-bold text-primary items-center"
              >
                {formType === "SIGN_IN" ? "Sign up" : "Sign in"} now
                <IconArrowNarrowRight />
              </Link>
            </p>
          </CardFooter>
        </Card>
      </form>
    </Form>
  );
};

export default AuthForm;
