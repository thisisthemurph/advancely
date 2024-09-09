export interface NavLinkProps {
  label: string;
  description: string;
  href: string;
  fullWidth?: boolean;
  showWhen?: "authenticated" | "unauthenticated" | "always";
}

export type NavLinks = NavLinkProps[];

export const navigationMenuItems: NavLinks = [
  {
    label: "Home",
    href: "/",
    description: "There's no place like home!",
  },
  {
    label: "About",
    href: "/",
    description: "Find out a little of who we are and what we can do for you.",
  },
  {
    label: "Pricing",
    href: "/",
    description:
      "Take a look at our competitive pricing models and see what best suits your needs.",
    fullWidth: true,
    showWhen: "unauthenticated",
  },
  {
    label: "Dashboard",
    href: "/dashboard",
    description: "Your dashboard contains all the things you need at the touch of a button.",
    fullWidth: true,
    showWhen: "authenticated",
  },
];
