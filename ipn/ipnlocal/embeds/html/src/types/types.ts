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

interface Service {
    Description: string;
    Port: number;
    Proto: string;
}

interface Profile {
    ID: string;
    LoginName: string;
    DisplayName: string;
    ProfilePicURL: string;
    Roles: object[];
}

interface Peer {
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

interface Base {
    factor: number; 
    suffix: string;
}

interface BytePrefix {
    short: string;
    full: string;
}
interface SSEMessage {
    length: number;
    type: string;
    data: any;
}

type BPrefix = Record<number, BytePrefix>

export type {
    Peer,
    AppData,
    Base,
    BPrefix,
    SSEMessage
}