"use client";

import { markTaskAsDone } from "@/api/actions/task";
import { Button } from "@/components/ui/button";
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
    <Button
      variant="outline"
      className="border-primary text-primary hover:text-primary"
      disabled={markingTaskAsDone}
      onClick={handleMarkTaskAsDone}
    >
      {markingTaskAsDone ? <Spinner /> : <IconChecks />}
      Mark as done
    </Button>
  );
};

export default MarkAsDoneBtn;
