"use client";

import Heading from "@/components/heading";
import StatusBadge from "@/components/StatusBadge";
import {
  IconCalendarWeek,
  IconClockPlus,
  IconLoader,
} from "@tabler/icons-react";
import { format } from "date-fns";
import React from "react";
import EditTaskBtn from "./EditBtn";
import MarkAsDoneBtn from "./MarkAsDoneBtn";
import DeleteBtn from "./DeleteBtn";

interface TaskProps {
  task: Task;
}

const Task = ({ task }: TaskProps) => {
  const { id, created_at, title, description, deadline, completed } = task;

  return (
    <>
      <Heading heading={title} />
      <div className="grid grid-cols-5 px-4 gap-y-6">
        <span className="flex gap-3 text-muted-foreground col-span-2 items-center">
          <IconClockPlus className="size-5" /> Created time
        </span>
        <span className="col-span-3">{format(created_at, "PPP HH:mm")}</span>
        <span className="flex gap-3 text-muted-foreground col-span-2 items-center">
          <IconLoader className="size-5" /> Status
        </span>
        <StatusBadge
          completed={completed}
          className="col-span-3"
          deadline={new Date(deadline)}
        />
        <span className="flex gap-3 text-muted-foreground col-span-2 items-center">
          <IconCalendarWeek className="size-5" /> Deadline
        </span>
        <span className="col-span-3">{format(deadline, "PPP HH:mm")}</span>
        {description && (
          <div className="col-span-full bg-muted p-6 rounded-2xl">
            <h3 className="text-lg font-semibold tracking-tight mb-4">
              Task Description
            </h3>
            <p className="text-muted-foreground text-sm whitespace-pre-line">
              {description}
            </p>
          </div>
        )}
        <div className="col-span-full mt-4 flex justify-between">
          <DeleteBtn id={id} />
          <div className="flex gap-4">
            {!completed && <MarkAsDoneBtn id={id} />}
            <EditTaskBtn id={id} />
          </div>
        </div>
      </div>
    </>
  );
};

export default Task;
