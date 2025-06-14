import React from "react";

interface HeadingProps {
  heading: string;
  Icon?: React.ComponentType<{ className: string }>;
  className?: string;
}

const Heading = ({ heading, Icon, className }: HeadingProps) => {
  return (
    <h1
      className={`mb-12 flex items-center gap-2 text-2xl font-bold ${className}`}
    >
      {Icon && <Icon className="size-7" />}
      <span className="line-clamp-2">{heading}</span>
    </h1>
  );
};

export default Heading;
