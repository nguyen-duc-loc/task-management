import Heading from "@/components/heading";
import { IconLayoutDashboardFilled } from "@tabler/icons-react";
import React, { Suspense } from "react";
import TaskContainer from "../_components/TaskContainer";
import Spinner from "@/components/ui/spinner";

const RootPage = ({
  searchParams,
}: {
  searchParams: Promise<{ [key: string]: string }>;
}) => {
  return (
    <div className="w-5xl">
      <Heading
        heading="Dashboard"
        Icon={IconLayoutDashboardFilled}
        className="mb-0!"
      />
      <p className="text-muted-foreground mb-12">
        Welcome to the task management
      </p>
      <Suspense
        fallback={
          <div>
            <Spinner className="mx-auto text-primary" />
          </div>
        }
      >
        <TaskContainer searchParams={searchParams} />
      </Suspense>
    </div>
  );
};

export default RootPage;
