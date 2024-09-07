import { cva, VariantProps } from "class-variance-authority";
import { cn } from "./ui/lib/utils.ts";

const logoVariants = cva("",
  {
    variants: {
      size: {
        sm: "w-6 h-6",
        md: "w-10 h-10",
        lg: "w-16 h-16",
        xl: "w-24 h-24"
      }
    },
    defaultVariants: {
      size: "md",
    },
  }
)

interface LogoProps extends Omit<React.HTMLAttributes<SVGElement>, "children">,
  VariantProps<typeof logoVariants> {}

const Logo = ({ className, size, ...props }: LogoProps) => (
  <svg
    id="Layer_1"
    xmlns="http://www.w3.org/2000/svg"
    xmlnsXlink="http://www.w3.org/1999/xlink"
    x="0px"
    y="0px"
    viewBox="0 0 492.481 492.481"
    xmlSpace="preserve"
    className={cn(logoVariants({ size, className }))}
    {...props}
  >
    <linearGradient
      id="SVGID_1_"
      gradientUnits="userSpaceOnUse"
      x1={-36.6002}
      y1={621.3422}
      x2={-17.2782}
      y2={547.7642}
      gradientTransform="matrix(7.8769 0 0 -7.8769 404.0846 4917.9966)"
    >
      <stop
        offset={0}
        style={{
          stopColor: "#29D3DA",
        }}
      />
      <stop
        offset={0.519}
        style={{
          stopColor: "#0077FF",
        }}
      />
      <stop
        offset={0.999}
        style={{
          stopColor: "#064093",
        }}
      />
      <stop
        offset={1}
        style={{
          stopColor: "#084698",
        }}
      />
    </linearGradient>
    <polygon
      style={{
        fill: "url(#SVGID_1_)",
      }}
      points="25.687,297.141 135.735,0 271.455,0 161.398,297.141 "
    />
    <linearGradient
      id="SVGID_2_"
      gradientUnits="userSpaceOnUse"
      x1={-27.0735}
      y1={620.7541}
      x2={-11.7045}
      y2={560.3241}
      gradientTransform="matrix(7.8769 0 0 -7.8769 404.0846 4917.9966)"
    >
      <stop
        offset={0.012}
        style={{
          stopColor: "#E0B386",
        }}
      />
      <stop
        offset={0.519}
        style={{
          stopColor: "#DA498C",
        }}
      />
      <stop
        offset={1}
        style={{
          stopColor: "#961484",
        }}
      />
    </linearGradient>
    <polygon
      style={{
        fill: "url(#SVGID_2_)",
      }}
      points="123.337,394.807 233.409,97.674 369.144,97.674 259.072,394.807 "
    />
    <linearGradient
      id="SVGID_3_"
      gradientUnits="userSpaceOnUse"
      x1={14.0324}
      y1={554.688}
      x2={-10.4176}
      y2={584.028}
      gradientTransform="matrix(7.8769 0 0 -7.8769 404.0846 4917.9966)"
    >
      <stop
        offset={0}
        style={{
          stopColor: "#29D3DA",
        }}
      />
      <stop
        offset={0.519}
        style={{
          stopColor: "#0077FF",
        }}
      />
      <stop
        offset={0.999}
        style={{
          stopColor: "#064093",
        }}
      />
      <stop
        offset={1}
        style={{
          stopColor: "#084698",
        }}
      />
    </linearGradient>
    <polygon
      style={{
        fill: "url(#SVGID_3_)",
      }}
      points="221.026,492.481 331.083,195.348 466.794,195.348 356.746,492.481 "
    />
  </svg>
);

export default Logo;
