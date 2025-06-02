import Link from "next/link";
import React from "react";

import ROUTES from "@/constants/routes";

import { Button } from "../ui/button";

interface ErrorPageProps {
  errorCode: number;
  message: string;
  description: string;
}

const ErrorPage = ({ errorCode, message, description }: ErrorPageProps) => {
  return (
    <div className="space-y-4 text-center">
      <span className="text-3xl font-semibold text-primary">{errorCode}</span>
      <p className="text-xl font-semibold">{message}</p>
      <p className="text-sm text-muted-foreground">{description}</p>
      <Button asChild variant="outline" className="!mt-10 rounded-lg">
        <Link href={ROUTES.dashboard} className="text-sm">
          Go to Dashboard
        </Link>
      </Button>
    </div>
  );
};

export default ErrorPage;
