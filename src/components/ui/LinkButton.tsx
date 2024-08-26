import { VariantProps } from "class-variance-authority";
import { buttonVariants } from "./button";

import { cn } from "./lib/utils";
import { Link, LinkProps } from "react-router-dom";
import { SheetClose } from "./sheet";

interface LinkButtonProps
  extends LinkProps,
    VariantProps<typeof buttonVariants> {}

function LinkButton({
  className,
  variant,
  size,
  children,
  ...props
}: LinkButtonProps) {
  return (
    <Link
      className={cn(buttonVariants({ variant, size, className }))}
      {...props}
    >
      {children}
    </Link>
  );
}

export function SheetCloseLinkButton(props: LinkButtonProps) {
  return (
    <SheetClose asChild>
      <LinkButton {...props} />
    </SheetClose>
  );
}

export default LinkButton;
