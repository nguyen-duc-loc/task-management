"use client";

import { deleteTask } from "@/api/actions/task";
import { DropdownMenuItem } from "@/components/ui/dropdown-menu";
import { IconTrash } from "@tabler/icons-react";
import React from "react";
import { toast } from "sonner";

interface DeleteBtnProps {
  id: string;
}

const DeleteBtn = ({ id }: DeleteBtnProps) => {
  const [deletingTask, startDeletingTask] = React.useTransition();

  const handleMarkTaskAsDone = () => {
    startDeletingTask(async () => {
      const result = await deleteTask(id);
      if (!result.success) {
        toast.error(`Failed to delete task. Try again later.`);
        return;
      }

      toast.success("Task has been deleted!");
    });
  };

  return (
    <DropdownMenuItem className="text-red-500 hover:text-red-500!" asChild>
      <button
        className="w-full"
        disabled={deletingTask}
        onClick={handleMarkTaskAsDone}
      >
        <IconTrash className="text-red-500" />
        Delete
      </button>
    </DropdownMenuItem>
  );
};

export default DeleteBtn;
