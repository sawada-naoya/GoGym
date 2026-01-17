"use client";

import { useFormContext } from "react-hook-form";

type Props = {
  name: string;
};

export const FormErrorMessage = ({ name }: Props) => {
  const {
    formState: { errors },
  } = useFormContext();

  const error = name.split(".").reduce((obj, key) => {
    return obj?.[key];
  }, errors as any);

  if (!error?.message) return null;

  return <p className="mt-1 text-xs text-red-600">{error.message as string}</p>;
};
