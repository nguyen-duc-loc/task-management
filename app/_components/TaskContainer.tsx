import { getTasks } from "@/api/data/task";
import React from "react";
import Task from "./Task";
import Pagination from "./Pagination";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import ROUTES from "@/constants/routes";
import { IconPlus } from "@tabler/icons-react";
import Search from "./Search";
import SelectStatus from "./SelectStatus";
import SelectDate from "./SelectDate";
import { isValid } from "date-fns";
import { Inbox } from "lucide-react";

const TaskContainer = async ({
  searchParams,
}: {
  searchParams: Promise<{ [key: string]: string }>;
}) => {
  const {
    page = "",
    completed = "",
    from = "",
    to = "",
    search = "",
  } = await searchParams;

  const currentPage = Math.max(Number(page) || 1, 1);
  const currentLimit = 12;

  const filterCompleted =
    completed === "true" ? true : completed === "false" ? false : undefined;

  const fromDate = isValid(new Date(from)) ? new Date(from) : undefined;
  const toDate = isValid(new Date(to)) ? new Date(to) : undefined;

  const { total, tasks } = await getTasks({
    limit: currentLimit,
    page: currentPage,
    completed: filterCompleted,
    fromDate,
    toDate,
    search,
  });

  return (
    <div className="space-y-6">
      <Button asChild>
        <Link href={ROUTES.newTask}>
          <IconPlus />
          Add New
        </Link>
      </Button>
      <div className="flex justify-between flex-wrap gap-4">
        <Search />
        <div className="flex flex-wrap gap-4">
          <SelectStatus />
          <SelectDate />
        </div>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        {tasks.length === 0 && (
          <div className="text-muted-foreground col-span-full flex flex-col items-center my-6">
            <Inbox className="size-8" />
            <span className="text-sm">No result</span>
          </div>
        )}
        {tasks.map((task) => (
          <Task task={task} key={task.id} />
        ))}
      </div>
      <Pagination
        total={total}
        currentLimit={currentLimit}
        currentPage={currentPage}
      />
    </div>
  );
};

export default TaskContainer;
