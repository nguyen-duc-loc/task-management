import React from "react";

interface BadgeProps {
  className?: string;
  completed: boolean;
  deadline: Date;
}

const StatusBadge = ({ className, completed, deadline }: BadgeProps) => {
  return (
    <div
      className={`w-fit rounded-lg px-3 py-1 text-center text-sm flex items-center gap-2 ${
        completed
          ? "bg-green-600/15"
          : deadline < new Date()
          ? "bg-red-600/15"
          : "bg-yellow-500/15"
      } ${className}`}
    >
      <div
        className={`w-2 h-2 rounded-full ${
          completed
            ? "bg-green-600"
            : deadline < new Date()
            ? "bg-red-600"
            : "bg-yellow-500"
        }`}
      />
      <span
        className={`${
          completed
            ? "text-green-600"
            : deadline < new Date()
            ? "text-red-600"
            : "text-yellow-500"
        }`}
      >
        {completed
          ? "Completed"
          : deadline < new Date()
          ? "Late"
          : "In progress"}
      </span>
    </div>
  );
};

export default StatusBadge;
