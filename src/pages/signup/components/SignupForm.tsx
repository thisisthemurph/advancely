import { z } from "zod";
import { useForm, UseFormReturn } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Input } from "../../../components/ui/input";
import { Button } from "../../../components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../../../components/ui/form";
import { useState } from "react";

const formSchema = z.object({
  name: z.string().min(4, {
    message: "Your company name should be at least 4 characters long",
  }),
  email: z.string().email(),
  password: z.string().min(6, {
    message: "Your password must be at least 6 characters long",
  }),
});

type FormInputs = z.infer<typeof formSchema>;

const CompanyNameInfo =
  "Set the name of your company, this is the name all of your employees will see themselves under when they sign in.";
const EmailInfo =
  "The email address used here will be considered a super admin account and will have permissions to create, read, update, and delete all content. Permissions can be tailored later to give other users different privilages.";

function SignupForm() {
  const [emailPlaceholder, setEmailPlaceholder] = useState(
    "your.name@company.com"
  );

  const form = useForm<FormInputs>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      email: "",
      password: "",
    },
  });

  function onSubmit(values: FormInputs) {
    console.log(values);
  }

  function onCompanyNameChange(form: UseFormReturn<FormInputs>) {
    const companyName = form.getValues("name");
    const email = `your.name@${companyName}.com`;
    setEmailPlaceholder(email.toLowerCase().replace(" ", ""));
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel infoText={EmailInfo}>Company name</FormLabel>
              <FormControl onChange={() => onCompanyNameChange(form)}>
                <Input
                  autoFocus={true}
                  placeholder="Your Company LTD"
                  {...field}
                />
              </FormControl>
              <FormMessage />
              <FormDescription className="hidden">
                {CompanyNameInfo}
              </FormDescription>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel infoText={EmailInfo}>
                <span>Email</span>
              </FormLabel>
              <FormControl>
                <Input type="email" placeholder={emailPlaceholder} {...field} />
              </FormControl>
              <FormMessage />
              <FormDescription className="hidden">{EmailInfo}</FormDescription>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <Input
                  type="password"
                  placeholder="* * * * * * * * * * * * * *"
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Sign up</Button>
      </form>
    </Form>
  );
}

export default SignupForm;
