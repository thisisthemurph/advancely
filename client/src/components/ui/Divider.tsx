interface Props {
  text: string;
}

function Divider({ text }: Props) {
  return (
    <div className="flex items-center my-8 w-full">
      <div className="flex-grow border-gray-300 border-t"></div>
      <span className="mx-4 font-semibold text-gray-500">{text}</span>
      <div className="flex-grow border-gray-300 border-t"></div>
    </div>
  );
}

export default Divider;
