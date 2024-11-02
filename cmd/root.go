/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beevik/ntp"
	"github.com/spf13/cobra"
)

// Variables to store DC IP Addr, original time, time shift duration, and change time
var originalTime, changeTime time.Time
var timeShift time.Duration
var dcIp string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "timeshift",
	Short: "Temporarily shift system time to match a domain controller",

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Store original time
		originalTime = time.Now()
		fmt.Println("Original system time stored: ", originalTime)

		// Get Windows Domain Controller time
		domainTime, err := getDomainNtpTime(dcIp)
		if err != nil {
			fmt.Println("Failed to get Domain Controller time:", err)
		}
		fmt.Println("Domain Controller Time:", domainTime)

		// Calculate and store time shift duration
		timeShift = domainTime.Sub(originalTime)
		fmt.Println("Time shift duration:", timeShift)

		// Set system time to domain time and record change time
		err = setSystemTime(domainTime)
		if err != nil {
			fmt.Println("Failed to set system time:", err)
			return
		}
		changeTime = time.Now()
		fmt.Println("System time set to domain controller time.")

		// Wait for user input to reset time
		fmt.Println("Press 'r' to reset time and exit...")
		waitForReset()

		// Calculate reset time
		elapsed := time.Since(changeTime)
		resetTime := originalTime.Add(elapsed)
		err = setSystemTime(resetTime)
		if err != nil {
			fmt.Println("Failed to reset system time:", err)
		} else {
			fmt.Println("System time reset to original:", resetTime)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringVarP(&dcIp, "dcIp", "d", "", "Domain Controller IP Address")
	rootCmd.MarkFlagRequired("dcIp")
}

func getDomainNtpTime(dcIp string) (time.Time, error) {
	dcTime, err := ntp.Time(dcIp)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to fetch DC Time: %v", err)
	}
	return dcTime, nil
}

func setSystemTime(newTime time.Time) error {
	tv := syscall.Timeval{
		Sec:  newTime.Unix(),
		Usec: int64(newTime.Nanosecond() / 1000),
	}
	return syscall.Settimeofday(&tv)
}

func waitForReset() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		fmt.Println("\nInterrupt signal received. Resetting time.")
	}()

	var input string
	for input != "r" {
		fmt.Scanln(&input)
		if input == "r" {
			fmt.Println("Resetting time...")
		} else {
			fmt.Println("Invalid input. Press 'r' to reset time.")
		}
	}
}
