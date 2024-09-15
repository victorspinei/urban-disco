import { Box, Flex, Button, useColorModeValue, useColorMode, Text, Container } from "@chakra-ui/react";
import { IoMoon } from "react-icons/io5";
import { LuSun } from "react-icons/lu";
import { FiGithub } from "react-icons/fi";
import { FaGithub } from "react-icons/fa";

export default function Navbar() {
	const { colorMode, toggleColorMode } = useColorMode();

	return (
		<Container maxW={"900px"}>
			<Box border={"2px"} borderColor={useColorModeValue("gray.200", "gray.800")} bg={useColorModeValue("white", "gray.700")} px={4} my={4} borderRadius={"5"}>
				<Flex h={16} alignItems={"center"} justifyContent={"space-between"}>
					{/* LEFT SIDE */}
					<Flex
						justifyContent={"center"}
						alignItems={"center"}
						gap={3}
						display={{ base: "none", sm: "flex" }}
					>
						<img src='/vinyl.png' alt='logo' width={40} height={40} />
						<Text fontSize={"xl"} fontWeight={500} >
							UrbanDisco 
						</Text>
					</Flex>

					{/* RIGHT SIDE */}
					<Flex alignItems={"center"} gap={6}>
                        <a href="https://github.com/victorspinei/urban-disco/">
							{colorMode === "light" ? <FiGithub size={20} /> : <FaGithub size={20} />}
                        </a>
						{/* Toggle Color Mode */}
						<Button onClick={toggleColorMode}>
							{colorMode === "light" ? <IoMoon /> : <LuSun size={20} />}
						</Button>
					</Flex>
				</Flex>
			</Box>
		</Container>
	);
}