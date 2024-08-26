import { cva, VariantProps } from "class-variance-authority";
import { cn } from "./lib/utils";

const logoVariants = cva(
  "font-semibold text-purple-400 hover:text-purple-600 hover:tracking-widest transition-all",
  {
    variants: {
      size: {
        default: "text-3xl",
        sm: "text-2xl",
        lg: "text-4xl",
      },
    },
    defaultVariants: {
      size: "default",
    },
  }
);

export interface LogoProps
  extends Omit<React.HTMLAttributes<HTMLHeadElement>, "children">,
    VariantProps<typeof logoVariants> {}

function Logo({ className, size, ...props }: LogoProps) {
  return (
    <h1 className={cn(logoVariants({ size, className }))} {...props}>
      advancely
    </h1>
  );
}

export default Logo;
