class AppData {
    Profile: Profile;
    DeviceName: string;
    NodeKey: string;
    StableID: string;
    Created: string;
    ServerURL: string;
    Version: string;
    Arch: string;
    OS: string;
    OSVersion: string;
    HostName: string;
    Name: string;
    Owner: string;
    Services: Service[];
    Health: string[];
    IPs: object[];
    IPv4: string;
    IPv6: string;
    Now: Date;
    TX: number;
    RX: number;
    Peers: Peer[];
}

class Service {
    Description: string;
    Port: number;
    Proto: string;
}

class Profile {
    ID: string;
    LoginName: string;
    DisplayName: string;
    ProfilePicURL: string;
    Roles: object[];
}

class Peer {
    HostName: string;
    ID: string;
    NodeKey: string;
    OS: string;
    Created: Date;
    DNSName: string;
    IPs: string[];
    IPv4: string;
    IPv6: string;
    RX: string;
    TX: string;
    RXb: number;
    TXb: number;
    Connection: string;
    ActAgo: string;
    LastSeen: Date;
    Unseen: boolean;
    CreatedDate: Date;
    CreatedTime: Date;
}



export type {
    Peer,
    AppData
}