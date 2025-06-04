"use client";

import { markTaskAsDone } from "@/api/actions/task";
import { Button } from "@/components/ui/button";
import { DropdownMenuItem } from "@/components/ui/dropdown-menu";
import Spinner from "@/components/ui/spinner";
import { IconChecks } from "@tabler/icons-react";
import React from "react";
import { toast } from "sonner";

interface MarkAsDoneBtnProps {
  id: string;
}

const MarkAsDoneBtn = ({ id }: MarkAsDoneBtnProps) => {
  const [markingTaskAsDone, startMarkingTaskAsDone] = React.useTransition();

  const handleMarkTaskAsDone = () => {
    startMarkingTaskAsDone(async () => {
      const result = await markTaskAsDone(id);
      if (!result.success) {
        toast.error(`Failed to mark task as done. Try again later.`);
        return;
      }

      toast.success("Task has been completed!");
    });
  };

  return (
    <DropdownMenuItem asChild>
      <button
        className="w-full"
        disabled={markingTaskAsDone}
        onClick={handleMarkTaskAsDone}
      >
        {markingTaskAsDone ? <Spinner /> : <IconChecks />}
        Mark as done
      </button>
    </DropdownMenuItem>
  );
};

export default MarkAsDoneBtn;
