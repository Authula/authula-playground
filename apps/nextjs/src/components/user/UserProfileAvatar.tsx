import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

type Props = {
  imageUrl: string;
  size?: "default" | "sm" | "lg";
};

export default function UserProfileAvatar({ imageUrl, size = "default" }: Props) {
  return (
    <Avatar size={size}>
      <AvatarImage src={imageUrl} alt="User Profile Image" />
      <AvatarFallback>
        <></>
      </AvatarFallback>
    </Avatar>
  );
}
