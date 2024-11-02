# timeshift
---

## Description
**timeshift** temporarily sets the local system time to match that of a specified domain controller in an active directory environment. This is helpful when the time skew between the attacking system and the target is too great for effective kerberos attacks.

## Installation

```bash
go install
```

## Usage

> [NOTE] 
> On Kali systems you may need to disable system ntp time like so:
> `sudo timedatectl set-ntp off`
> before running **timeshift**.
> To enable it after you're done and have reset the time just run:
> `sudo timedatectl set-ntp on`

To set time to that of the DC:

```bash
timeshift -d <DC_IP>
```